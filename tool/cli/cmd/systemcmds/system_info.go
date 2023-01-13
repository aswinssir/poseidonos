package systemcmds

import (
	"cli/cmd/displaymgr"
	"cli/cmd/globals"
	"cli/cmd/grpcmgr"
	"cli/cmd/otelmgr"
	"fmt"

	"github.com/spf13/cobra"
)

var SystemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about PoseidonOS and server hardware.",
	Long: `
Display information about PoseidonOS and server hardware such as BIOS, baseboard, and processors.

Syntax:
	poseidonos-cli system info
          `,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Start OpenTelemetry trace
		m := otelmgr.GetOtelManagerInstance()
		defer m.Shutdown()
		t := otelmgr.NewTracer()
		t.SetTrace(m.GetRootContext(), globals.SYSTEM_CMD_APP_NAME, globals.SYSTEM_INFO_FUNC_NAME)
		defer t.Release()

		posMgr, err := grpcmgr.GetPOSManager()
		if err != nil {
			fmt.Printf("failed to connect to POS: %v", err)
			return err
		}
		res, req, gRpcErr := posMgr.StartTelemetry()

		printReqErr := displaymgr.PrintProtoReqJson(req)
		if printReqErr != nil {
			fmt.Printf("failed to marshal the protobuf request: %v", printReqErr)
			t.RecordError(printReqErr)
			return printReqErr
		}

		if gRpcErr != nil {
			globals.PrintErrMsg(gRpcErr)
			t.RecordError(gRpcErr)
			return gRpcErr
		}

		printResErr := displaymgr.PrintProtoResponse(req.Command, res)
		if printResErr != nil {
			fmt.Printf("failed to print the response: %v", printResErr)
			t.RecordError(printResErr)
			return printResErr
		}

		return nil
	},
}

func init() {
}
