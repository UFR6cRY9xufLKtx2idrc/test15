/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package tcp_server

import (
	"fmt"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/coremain"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/pkg/server/http_handler"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/pkg/utils"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/server/server_utils"
	"golang.org/x/net/http2"
	"net/http"
	"time"
)

const PluginType = "http_server"

func init() {
	coremain.RegNewPluginFunc(PluginType, Init, func() any { return new(Args) })
}

type Args struct {
	Entries []struct {
		Exec string `yaml:"exec"`
		Path string `yaml:"path"`
	} `yaml:"entries"`
	Listen      string `yaml:"listen"`
	SrcIPHeader string `yaml:"src_ip_header"`
	Cert        string `yaml:"cert"`
	Key         string `yaml:"key"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

func (a *Args) init() {
	utils.SetDefaultNum(&a.IdleTimeout, 30)
}

type HttpServer struct {
	args *Args

	server *http.Server
}

func (s *HttpServer) Close() error {
	return s.server.Close()
}

func Init(bp *coremain.BP, args any) (any, error) {
	return StartServer(bp, args.(*Args))
}

func StartServer(bp *coremain.BP, args *Args) (*HttpServer, error) {
	mux := http.NewServeMux()
	for _, entry := range args.Entries {
		dh, err := server_utils.NewHandler(bp, entry.Exec)
		if err != nil {
			return nil, fmt.Errorf("failed to init dns handler, %w", err)
		}
		hhOpts := http_handler.HandlerOpts{
			DNSHandler:  dh,
			SrcIPHeader: args.SrcIPHeader,
			Logger:      bp.L(),
		}
		hh := http_handler.NewHandler(hhOpts)
		mux.Handle(entry.Path, hh)
	}

	hs := &http.Server{
		Addr:           args.Listen,
		Handler:        mux,
		ReadTimeout:    time.Second,
		IdleTimeout:    time.Duration(args.IdleTimeout) * time.Second,
		MaxHeaderBytes: 512,
	}
	if err := http2.ConfigureServer(hs, &http2.Server{
		MaxReadFrameSize:             16 * 1024,
		IdleTimeout:                  time.Duration(args.IdleTimeout) * time.Second,
		MaxUploadBufferPerConnection: 65535,
		MaxUploadBufferPerStream:     65535,
	}); err != nil {
		return nil, fmt.Errorf("failed to setup http2 server, %w", err)
	}

	go func() {
		var err error
		if len(args.Key)+len(args.Cert) > 0 {
			err = hs.ListenAndServeTLS(args.Cert, args.Key)
		} else {
			err = hs.ListenAndServe()
		}
		bp.M().GetSafeClose().SendCloseSignal(err)
	}()
	return &HttpServer{
		args:   args,
		server: hs,
	}, nil
}
