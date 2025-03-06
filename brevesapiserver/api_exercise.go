/*
 * Breves API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package brevesapiserver

import (
	"encoding/json"
	"fmt"
	"github.com/ZolaraProject/breves-api/models"
	brevesVaultService "github.com/ZolaraProject/breves-vault-service/brevesvaultrpc"
	grpctoken "github.com/ZolaraProject/library/grpctoken"
	logger "github.com/ZolaraProject/library/logger"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"path"
	"strconv"
)

func GetUserVideos(w http.ResponseWriter, r *http.Request) {
	ctx, grpcToken := grpctoken.CreateContextFromHeader(r, JwtSecretKey)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		logger.Err(grpcToken, "failed to get metadata from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get metadata from context")
		return
	}

	if len(md.Get("zolara-user-id")) == 0 {
		logger.Err(grpcToken, "failed to get user id from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get user id from context")
		return
	}

	userId, err := strconv.ParseInt(md.Get("zolara-user-id")[0], 10, 64)
	if err != nil {
		logger.Err(grpcToken, "Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err))
		return
	}

	// Create gRPC client
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", PkiVaultServiceHost, PkiVaultServicePort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Err(grpcToken, "RegisterUser could not establish gRPC connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("RegisterUser could not establish gRPC connection: %v", err))
		return
	}
	defer conn.Close()
	client := brevesVaultService.NewBrevesVaultServiceClient(conn)

	res, err := client.GetUserVideos(ctx, &brevesVaultService.UserVideoRequest{UserId: userId})
	if err != nil {
		logger.Err(grpcToken, "failed to get user videos: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get user videos")
		return
	}

	userVideos := []models.UserVideoInList{}
	for _, video := range res.GetUserVideos() {
		userVideos = append(userVideos, models.UserVideoInList{
			Title:    video.Title,
			Subtitle: video.Subtitle,
			Likes:    video.Likes,
			Language: fmt.Sprintf("%s", video.Language),
			Level:    fmt.Sprintf("%s", video.Level),
			Action:   video.Action,
			VideoUrl: video.VideoUrl,
			VideoId:  video.VideoId,
		})
	}

	response, err := json.Marshal(&models.UserVideoList{
		UserVideos: userVideos,
		Total:      res.GetTotal(),
	})
	if err != nil {
		logger.Err(grpcToken, "failed to marshal response: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func LikeVideos(w http.ResponseWriter, r *http.Request) {
	ctx, grpcToken := grpctoken.CreateContextFromHeader(r, JwtSecretKey)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		logger.Err(grpcToken, "failed to get metadata from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get metadata from context")
		return
	}

	if len(md.Get("zolara-user-id")) == 0 {
		logger.Err(grpcToken, "failed to get user id from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get user id from context")
		return
	}

	userId, err := strconv.ParseInt(md.Get("zolara-user-id")[0], 10, 64)
	if err != nil {
		logger.Err(grpcToken, "Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err))
		return
	}

	var videos models.VideoLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&videos); err != nil {
		logger.Err(grpcToken, "failed to decode request body: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("failed to decode request body: %s", err))
		return
	}

	lifeReq := []*brevesVaultService.VideoLikes{}
	for _, video := range videos.Videos {
		videoId, err := strconv.ParseInt(path.Base(video), 10, 64)
		if err != nil {
			logger.Err(grpcToken, "failed to parse video id: %s", err)
			continue
		}

		lifeReq = append(lifeReq, &brevesVaultService.VideoLikes{
			Id: videoId,
		})
	}

	// Create gRPC client
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", PkiVaultServiceHost, PkiVaultServicePort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Err(grpcToken, "RegisterUser could not establish gRPC connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("RegisterUser could not establish gRPC connection: %v", err))
		return
	}
	defer conn.Close()
	client := brevesVaultService.NewBrevesVaultServiceClient(conn)

	res, err := client.LikeVideo(ctx, &brevesVaultService.LikeVideoRequest{
		UserId:     userId,
		VideoLikes: lifeReq,
	})
	if err != nil {
		logger.Err(grpcToken, "failed to like videos: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("failed to like videos: %s", err))
		return
	}

	response, err := json.Marshal(&models.UserCreatedResponse{
		Message: res.GetMessage(),
	})
	if err != nil {
		logger.Err(grpcToken, "failed to marshal response: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateUserVideo(w http.ResponseWriter, r *http.Request) {
	ctx, grpcToken := grpctoken.CreateContextFromHeader(r, JwtSecretKey)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		logger.Err(grpcToken, "failed to get metadata from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get metadata from context")
		return
	}

	if len(md.Get("zolara-user-id")) == 0 {
		logger.Err(grpcToken, "failed to get user id from context")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, "failed to get user id from context")
		return
	}

	vars := mux.Vars(r)
	i := vars["videoId"]
	if len(i) == 0 {
		logger.Err(grpcToken, "videoId is required")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeStandardResponse(r, w, grpcToken, "id is required")
		return
	}

	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		logger.Err(grpcToken, "failed to parse id into int64: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeStandardResponse(r, w, grpcToken, "failed to parse id into int64")
		return
	}

	userId, err := strconv.ParseInt(md.Get("zolara-user-id")[0], 10, 64)
	if err != nil {
		logger.Err(grpcToken, "Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("Failed to convert user id(%s) to Int : %s", md.Get("zolara-user-id")[0], err))
		return
	}

	// Create gRPC client
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", PkiVaultServiceHost, PkiVaultServicePort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Err(grpcToken, "RegisterUser could not establish gRPC connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("RegisterUser could not establish gRPC connection: %v", err))
		return
	}
	defer conn.Close()
	client := brevesVaultService.NewBrevesVaultServiceClient(conn)

	res, err := client.UpdateUserVideo(ctx, &brevesVaultService.UpdateUserVideoRequest{
		UserId:  userId,
		VideoId: id,
	})
	if err != nil {
		logger.Err(grpcToken, "failed to update user's understanding of video: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeStandardResponse(r, w, grpcToken, fmt.Sprintf("failed to update user's understanding of video"))
		return
	}

	var response []byte
	if len(res.GetUserVideos()) == 0 {
		response, err = json.Marshal(&models.UserVideoInList{})
		if err != nil {
			logger.Err(grpcToken, "failed to marshal response: %s", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			writeStandardResponse(r, w, grpcToken, fmt.Sprintf("failed to marshal response"))
			return
		}
	} else {
		userVideos := []models.UserVideoInList{}
		for _, video := range res.GetUserVideos() {
			userVideos = append(userVideos, models.UserVideoInList{
				Title:    video.Title,
				Subtitle: video.Subtitle,
				Likes:    video.Likes,
				Language: fmt.Sprintf("%s", video.Language),
				Level:    fmt.Sprintf("%s", video.Level),
				Action:   video.Action,
				VideoUrl: video.VideoUrl,
				VideoId:  video.VideoId,
			})
		}

		response, err = json.Marshal(&models.UserVideoList{
			UserVideos: userVideos,
			Total:      res.GetTotal(),
		})
		if err != nil {
			logger.Err(grpcToken, "failed to marshal response: %s", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			writeStandardResponse(r, w, grpcToken, fmt.Sprintf("failed to marshal response"))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
