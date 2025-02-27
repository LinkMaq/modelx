package model

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"kubegems.io/modelx/cmd/modelx/repo"
	"kubegems.io/modelx/pkg/client"
)

func NewLoginCmd() *cobra.Command {
	token := ""
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to a modelx repository",
		Example: `
  modelx repo add local http://localhost:8080
  modelx login repo local --token <token>
		`,
		SilenceUsage: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return repo.CompleteRegistry(toComplete)
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
			defer cancel()
			if len(args) == 0 {
				return errors.New("at least one argument is required")
			}
			if token == "" {
				fmt.Print("Token: ")
				fmt.Scanln(&token)
			}
			return LoginModelx(ctx, args[0], token)
		},
	}
	cmd.Flags().StringVarP(&token, "token", "t", "", "token")
	return cmd
}

func LoginModelx(ctx context.Context, reponame string, token string) error {
	repoDetails, err := repo.DefaultRepoManager.Get(reponame)
	if err != nil {
		return err
	}
	repoDetails.Token = token
	if err := repoDetails.Client().Ping(ctx); err != nil {
		return err
	}
	fmt.Printf("Login successful for %s\n", reponame)
	return repo.DefaultRepoManager.Set(repoDetails)
}

func Ping(ctx context.Context, repo string, token string) error {
	token = "Bearer " + token
	return client.NewClient(repo, token).Ping(ctx)
}
