package webdavrpc

import (
	"log"

	"github.com/zv0n/webdav-proxy/configuration"
	"github.com/zv0n/webdav-proxy/webdav"
	"golang.org/x/net/context"
)

type Server struct {
	Config *configuration.Configuration
}

func (s *Server) MountWebdav(ctx context.Context, request *MountWebdavRequest) (*MountWebdavResponse, error) {
	log.Printf("Received mount request: url => \"%s\"; dir => \"%s\"", request.Url, request.Dir)
	err := webdav.Mount(webdav.MountInput{
		URL:        request.Url,
		Dir:        request.Dir,
		User:       request.User,
		Password:   request.Password,
		ConfigName: request.ConfigName,
		TargetPath: request.Target,
		UID:        request.Uid,
		GID:        request.Gid,
	}, s.Config)
	if err != nil {
		return &MountWebdavResponse{Output: err.Error()}, err
	}
	return &MountWebdavResponse{Output: "Success"}, nil
}

func (s *Server) UmountWebdav(ctx context.Context, request *UmountWebdavRequest) (*UmountWebdavResponse, error) {
	log.Printf("Received umount request from client: target => %s", request.MountTarget)
	err := webdav.Umount(request.MountTarget, request.ConfigName)
	if err != nil {
		return &UmountWebdavResponse{Output: err.Error()}, err
	}
	return &UmountWebdavResponse{Output: "Success"}, nil
}
