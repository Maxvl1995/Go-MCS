package gapi

import (
	"Go-MCS/models"
	"Go-MCS/pb"
	"Go-MCS/utils"
	"context"
	"strings"

	"github.com/thanhpk/randstr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.GenericResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.SignUpInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	newUser, err := server.authService.SignUpUser(&user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	server.userService.UpdateUserById(newUser.ID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[0]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       server.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	// err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())

	// }

	utils.SendEmail(newUser, &emailData)

	message := "We sent an email with a verification code to " + newUser.Email

	res := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return res, nil
}
