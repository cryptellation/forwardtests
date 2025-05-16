//go:build integration
// +build integration

package sql

import (
	"context"
	"testing"

	"github.com/cryptellation/dbmigrator"
	"github.com/cryptellation/forwardtests/configs"
	"github.com/cryptellation/forwardtests/configs/sql/down"
	"github.com/cryptellation/forwardtests/configs/sql/up"
	"github.com/cryptellation/forwardtests/svc/db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestForwardtestsSuite(t *testing.T) {
	suite.Run(t, new(ForwardtestsSuite))
}

type ForwardtestsSuite struct {
	db.ForwardtestSuite
}

func (suite *ForwardtestsSuite) SetupSuite() {
	act, err := New(context.Background(), viper.GetString(configs.EnvSQLDSN))
	suite.Require().NoError(err)

	mig, err := dbmigrator.NewMigrator(context.Background(), act.db, up.Migrations, down.Migrations, nil)
	suite.Require().NoError(err)
	suite.Require().NoError(mig.MigrateToLatest(context.Background()))

	suite.DB = act
}

func (suite *ForwardtestsSuite) SetupTest() {
	db := suite.DB.(*Activities)
	suite.Require().NoError(db.Reset(context.Background()))
}
