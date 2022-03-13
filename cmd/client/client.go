package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/adinandradrs/cadence-sample/cmd"
	"github.com/adinandradrs/cadence-sample/docs"
	"github.com/adinandradrs/cadence-sample/model"
	"github.com/adinandradrs/cadence-sample/workflows"
	"github.com/gin-gonic/gin"
	ginswag "github.com/swaggo/gin-swagger"
	swagfile "github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

// @title Sandbox API - Cadence Client Service
// @version 0.1
// @description This is a sandbox APIs for Cadence Client
// @termsOfService https://adinandra.dharmasurya.id

// @contact.name Adinandra Dharmasurya
// @contact.url https://adinandra.dharmasurya.id
// @contact.email adinandra.dharmasurya@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
func main() {
	logger := cmd.Logger()
	cmd.Environment(logger)
	cadclient := cmd.CadenceClient(os.Getenv("CADENCE_HOST"), os.Getenv("CADENCE_DOMAIN"))
	regctrl := newRegisterController(cadclient, logger)
	route := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	route.GET("/swagger/*any", ginswag.WrapHandler(swagfile.Handler))
	route.POST("/v1/register", regctrl.ExecuteRegister)
	route.POST("/v1/decision", regctrl.ExecuteDecision)

	err := route.Run(os.Getenv("APP_PORT"))
	if err != nil {
		panic("client error - failed to running")
	}
}

type register struct {
	cadclient client.Client
	logger    *zap.Logger
}

type registerController interface {
	ExecuteRegister(ctx *gin.Context)
	ExecuteDecision(ctx *gin.Context)
}

func newRegisterController(c client.Client, l *zap.Logger) registerController {
	return &register{
		cadclient: c,
		logger:    l,
	}
}

// @BasePath /api/v1
// @Tags Register APIs
// API Register
// @Summary API Register Form
// @Schemes
// @Description API to register an user
// @Accept json
// @Produce json
// @Param payload body model.RegisterRequest true "Register Request Payload"
// @Success 200 "Request successfuly processed"
// @Failure 400 "Got a bad payload"
// @Failure 401 "Unauthorized access"
// @Failure 403 "Forbidden access"
// @Failure 404 "Data is not found"
// @Failure 500 "Something went wrong"
// @Router /v1/register [post]
func (r *register) ExecuteRegister(ctx *gin.Context) {
	req := model.RegisterRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if !ctx.IsAborted() {
		opt := client.StartWorkflowOptions{
			ID:                              "Register - " + req.Email,
			TaskList:                        cmd.TaskList,
			ExecutionStartToCloseTimeout:    time.Minute * 5,
			DecisionTaskStartToCloseTimeout: time.Minute * 5,
			WorkflowIDReusePolicy:           client.WorkflowIDReusePolicyAllowDuplicate,
		}
		result, err := r.cadclient.StartWorkflow(context.Background(), opt, workflows.RegisterFlow, req)
		if err != nil {
			r.logger.Error("execute workflow register got an error", zap.Error(err), zap.String("", req.Email))
			ctx.AbortWithStatus(http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, model.WorkflowResponse{
				WorkflowId:    result.ID,
				WorkflowRunId: result.RunID,
			})
		}
	}
}

// @BasePath /api/v1
// @Tags Register APIs
// API Register Decision
// @Summary API Register Decision
// @Schemes
// @Description API to give a decision for registration form
// @Accept json
// @Produce json
// @Param payload body model.ApprovalRequest true "Approval Request Payload"
// @Success 200 "Request successfuly processed"
// @Failure 400 "Got a bad payload"
// @Failure 401 "Unauthorized access"
// @Failure 403 "Forbidden access"
// @Failure 404 "Data is not found"
// @Failure 500 "Something went wrong"
// @Router /v1/decision [post]
func (r *register) ExecuteDecision(ctx *gin.Context) {
	req := model.ApprovalRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if !ctx.IsAborted() {
		err := r.cadclient.CompleteActivityByID(context.Background(), os.Getenv("CADENCE_DOMAIN"),
			req.WorkflowId,
			req.WorkflowRunId,
			req.ActivityId,
			req, nil,
		)
		if err != nil {
			r.logger.Error("complete activity error", zap.Error(err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, "Completed")
		}
	}

}
