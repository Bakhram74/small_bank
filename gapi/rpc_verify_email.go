package gapi

import (
	"context"
	db "github.com/Bakhram74/small_bank/db/sqlc"
	"github.com/Bakhram74/small_bank/pb"
	"github.com/Bakhram74/small_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	verifyEmailTx, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}
	rps := &pb.VerifyEmailResponse{IsVerified: verifyEmailTx.User.IsEmailVerified}

	return rps, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidateEmailSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
