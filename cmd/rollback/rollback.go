package rollback

import (
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rollbackOptions struct {
	Message string
	// if customer want do rollback on failed argocd sync
	NeedRollback bool
}

var Cmd = &cobra.Command{
	Use:   "rollback [app]",
	Short: "Handle rollback case",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cf := codefresh.New(&codefresh.ClientOptions{
			Token: context.PluginCodefreshCredentials.Token,
			Host:  context.PluginCodefreshCredentials.Host,
		})

		if rollbackOptions.Message != "" && rollbackOptions.NeedRollback {
			log.Println(fmt.Sprintf("Start do rollback because message is \"%s\"", rollbackOptions.Message))

			rollbackResult, _ := cf.RollbackToStable(name, codefresh.Rollback{
				ContextName:     context.PluginCodefreshCredentials.Integration,
				ApplicationName: name,
			})

			if rollbackResult != nil {
				log.Println(fmt.Sprintf("Run rollback process pipeline, build link https://g.codefresh.io/build/%s", rollbackResult.BuildId))
			}

			return errors.New("running rollback flow")
		}

		log.Println("Skip rollback execution")

		return nil
	},
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&rollbackOptions.Message, "message", "", "Error message from sync execution")
	f.BoolVar(&rollbackOptions.NeedRollback, "needRollback", false, "Execute rollback if sync is failed")
}
