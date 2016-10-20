package main

import (
	"github.com/dtan4/ec2c/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"cancel": func() (cli.Command, error) {
			return &command.CancelCommand{
				Meta: *meta,
			}, nil
		},
		"launch": func() (cli.Command, error) {
			return &command.LaunchCommand{
				Meta: *meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"list-requests": func() (cli.Command, error) {
			return &command.ListRequestsCommand{
				Meta: *meta,
			}, nil
		},
		"request": func() (cli.Command, error) {
			return &command.RequestCommand{
				Meta: *meta,
			}, nil
		},
		"tag": func() (cli.Command, error) {
			return &command.TagCommand{
				Meta: *meta,
			}, nil
		},
		"untag": func() (cli.Command, error) {
			return &command.UntagCommand{
				Meta: *meta,
			}, nil
		},
		"terminate": func() (cli.Command, error) {
			return &command.TerminateCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: Revision,
				Name:     Name,
			}, nil
		},
	}
}
