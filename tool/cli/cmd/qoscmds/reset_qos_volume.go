package qoscmds

import (
	"encoding/json"
	"pnconnector/src/log"
	"strings"

	"cli/cmd/displaymgr"
	"cli/cmd/globals"
	"cli/cmd/messages"
	"cli/cmd/socketmgr"

	"github.com/spf13/cobra"
)

var VolumeResetCmd = &cobra.Command{
	Use:   "reset [flags]",
	Short: "Reset qos policy for a volume(s) of PoseidonOS.",
	Long: `Reset qos policy for a volume of PoseidonOS.

Syntax: 
	poseidonos-cli qos reset (--volume-name | -v) VolumeName (--array-name | -a) ArrayName .

Example: 
	poseidonos-cli qos reset --volume-name Volume0 --array-name Array0
          `,

	Run: func(cmd *cobra.Command, args []string) {

		var command = "RESETQOSVOLUMEPOLICY"

		volumeResetReq := formVolumeResetReq()
		reqJSON, err := json.Marshal(volumeResetReq)
		if err != nil {
			log.Debug("error:", err)
		}

		displaymgr.PrintRequest(string(reqJSON))

		// Do not send request to server and print response when testing request build.
		if !(globals.IsTestingReqBld) {
			socketmgr.Connect()

			resJSON, err := socketmgr.SendReqAndReceiveRes(string(reqJSON))
			if err != nil {
				log.Debug("error:", err)
				return
			}

			socketmgr.Close()

			displaymgr.PrintResponse(command, resJSON, globals.IsDebug, globals.IsJSONRes, globals.DisplayUnit)
		}
	},
}

func formVolumeResetReq() messages.Request {

	volumeNameListSlice := strings.Split(volumeReset_volumeNameList, ",")
	var volumeNames []messages.VolumeNameList
	for _, str := range volumeNameListSlice {
		var volumeNameList messages.VolumeNameList // Single device name that is splitted
		volumeNameList.VOLUMENAME = str
		volumeNames = append(volumeNames, volumeNameList)
	}

	volumeResetParam := messages.VolumePolicyParam{
		VOLUMENAME: volumeNames,
		ARRAYNAME:  volumeReset_arrayName,
	}

	volumeResetReq := messages.Request{
		RID:     "fromCLI",
		COMMAND: "RESETQOSVOLUMEPOLICY",
		PARAM:   volumeResetParam,
	}

	return volumeResetReq
}

// Note (mj): In Go-lang, variables are shared among files in a package.
// To remove conflicts between variables in different files of the same package,
// we use the following naming rule: filename_variablename. We can replace this if there is a better way.
var volumeReset_volumeNameList = ""
var volumeReset_arrayName = ""

func init() {
	VolumeResetCmd.Flags().StringVarP(&volumeReset_volumeNameList, "volume-name", "v", "", "A comma-seperated names of volumes to set qos policy for")
	VolumeResetCmd.MarkFlagRequired("volume-name")

	VolumeResetCmd.Flags().StringVarP(&volumeReset_arrayName, "array-name", "a", "", "Name of the array where the volume is created from")
	VolumeResetCmd.MarkFlagRequired("array-name")
}
