package arraycmds

import (
	"cli/cmd/displaymgr"
	"cli/cmd/globals"
	"cli/cmd/grpcmgr"
	"fmt"
	pb "kouros/api"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var ReplaceArrayDeviceCmd = &cobra.Command{
	Use:   "replace [flags]",
	Short: "Replace a data device with an available spare device in array.",
	Long: `
Replace a data device with an available spare device in array. Use this command when you expect
a possible problem of a data device. If there is no available spare device, this command will fail.

Syntax:
	poseidonos-cli array replace (--data-device | -d) DeviceName (--array-name | -a) ArrayName

Example: 
	poseidonos-cli array replace --data-device nvme5 --array-name array0
          `,
	RunE: func(cmd *cobra.Command, args []string) error {

		reqParam, buildErr := buildReplaceArrayDeviceReqParam()
		if buildErr != nil {
			fmt.Printf("failed to build request: %v", buildErr)
			return buildErr
		}

		posMgr, err := grpcmgr.GetPOSManager()
		if err != nil {
			fmt.Printf("failed to connect to POS: %v", err)
			return err
		}
		res, req, gRpcErr := posMgr.ReplaceArrayDevice(reqParam)

		reqJson, err := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}.Marshal(req)
		if err != nil {
			fmt.Printf("failed to marshal the protobuf request: %v", err)
			return err
		}
		displaymgr.PrintRequest(string(reqJson))

		if gRpcErr != nil {
			globals.PrintErrMsg(gRpcErr)
			return gRpcErr
		}

		printErr := displaymgr.PrintProtoResponse(req.Command, res)
		if printErr != nil {
			fmt.Printf("failed to print the response: %v", printErr)
			return printErr
		}

		return nil
	},
}

func buildReplaceArrayDeviceReqParam() (*pb.ReplaceArrayDeviceRequest_Param, error) {
	param := &pb.ReplaceArrayDeviceRequest_Param{Array: replace_array_device_arrayName, Device: replace_array_device_dataDev}

	return param, nil
}

// Note (mj): In Go-lang, variables are shared among files in a package.
// To remove conflicts between variables in different files of the same package,
// we use the following naming rule: filename_variablename. We can replace this if there is a better way.
var replace_array_device_arrayName = ""
var replace_array_device_dataDev = ""

func init() {
	ReplaceArrayDeviceCmd.Flags().StringVarP(&replace_array_device_arrayName,
		"array-name", "a", "",
		"The name of the array of the data and spare devices.")
	ReplaceArrayDeviceCmd.MarkFlagRequired("array-name")

	ReplaceArrayDeviceCmd.Flags().StringVarP(&replace_array_device_dataDev,
		"data-device", "d", "",
		"The name of the device to be replaced with.")
	ReplaceArrayDeviceCmd.MarkFlagRequired("data-device")
}
