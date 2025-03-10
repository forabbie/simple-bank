package validator

import (
	"github.com/forabbie/vank-app/models"
	"github.com/forabbie/vank-app/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateLoginUserRequest(req *models.LoginUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if err := ValidateUsername(req.Username); err != nil {
		violations = append(violations, util.CreateFieldViolation("username", err))
	}

	if err := ValidatePassword(req.Password); err != nil {
		violations = append(violations, util.CreateFieldViolation("password", err))
	}

	return violations
}

func ValidateCreateUserRequest(req *models.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if err := ValidateUsername(req.Username); err != nil {
		violations = append(violations, util.CreateFieldViolation("username", err))
	}

	if err := ValidatePassword(req.Password); err != nil {
		violations = append(violations, util.CreateFieldViolation("password", err))
	}

	if err := ValidateFullName(req.FullName); err != nil {
		violations = append(violations, util.CreateFieldViolation("full_name", err))
	}

	if err := ValidateEmail(req.Email); err != nil {
		violations = append(violations, util.CreateFieldViolation("email", err))
	}

	return violations
}
