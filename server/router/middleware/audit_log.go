package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
)

var auditRespPool = sync.Pool{
	New: func() any {
		return &auditResponseWriter{
			ResponseWriter: nil,
			body:           &bytes.Buffer{},
		}
	},
}

type auditResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r auditResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		writer := auditRespPool.Get().(*auditResponseWriter)
		defer auditRespPool.Put(writer)
		writer.ResponseWriter = c.Writer
		c.Writer = writer

		c.Next()

		respBody := writer.body.Bytes()
		writer.body.Reset()

		claimsVal, exists := c.Get("claims")
		if !exists {
			return
		}
		claims := claimsVal.(jwt.Claims)

		path := c.Request.RequestURI
		trimmed := strings.TrimPrefix(path, "/api/v1/admin/")
		idx := strings.Index(trimmed, "?")
		if idx > 0 {
			trimmed = trimmed[:idx]
		}
		parts := strings.Split(trimmed, "/")
		moduleName := ""
		action := ""
		if len(parts) >= 2 {
			moduleName = strings.Join(parts[:len(parts)-1], "/")
			action = parts[len(parts)-1]
		} else if len(parts) == 1 {
			moduleName = parts[0]
			action = parts[0]
		}

		respCode := 0
		var respMap map[string]any
		if json.Unmarshal(respBody, &respMap) == nil {
			if code, ok := respMap["code"]; ok {
				if v, ok := code.(float64); ok {
					respCode = int(v)
				}
			}
		}

		reqStr := string(reqBody)
		if len(reqStr) > 4096 {
			reqStr = reqStr[:4096]
		}

		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		ms := time.Since(start).Milliseconds()

		go func() {
			db, _, _ := repository.Get("")
			gdb := db.UnderlyingDB()
			_ = gdb.Table("system_audit_log").Create(&model.SystemAuditLog{
				UserCode:  claims.Code,
				Module:    moduleName,
				Action:    action,
				Path:      path,
				ReqBody:   reqStr,
				RespCode:  respCode,
				IP:        ip,
				UserAgent: ua,
				Ms:        ms,
			}).Error
		}()
	}
}
