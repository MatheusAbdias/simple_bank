package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/MatheusAbdias/go_simple_bank/config"
	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
