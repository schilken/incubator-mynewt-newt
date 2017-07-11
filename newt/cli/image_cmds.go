/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"mynewt.apache.org/newt/newt/builder"
	"mynewt.apache.org/newt/newt/image"
	"mynewt.apache.org/newt/newt/newtutil"
	"mynewt.apache.org/newt/util"
)

func createImageRunCmd(cmd *cobra.Command, args []string) {
	var keyId uint8
	var keystr string

	if len(args) < 2 {
		NewtUsage(cmd, util.NewNewtError("Must specify target and version"))
	}

	TryGetProject()

	targetName := args[0]
	t := ResolveTarget(targetName)
	if t == nil {
		NewtUsage(cmd, util.NewNewtError("Invalid target name: "+targetName))
	}

	version := args[1]

	if len(args) > 2 {
		if len(args) > 3 {
			keyId64, err := strconv.ParseUint(args[3], 10, 8)
			if err != nil {
				NewtUsage(cmd,
					util.NewNewtError("Key ID must be between 0-255"))
			}
			keyId = uint8(keyId64)
		}
		keystr = args[2]
	}

	b, err := builder.NewTargetBuilder(t)
	if err != nil {
		NewtUsage(nil, err)
	}

	if _, _, err := b.CreateImages(version, keystr, keyId); err != nil {
		NewtUsage(nil, err)
		return
	}
}

func resignImageRunCmd(cmd *cobra.Command, args []string) {
	var keyId uint8
	var keystr string

	if len(args) < 1 {
		NewtUsage(cmd, util.NewNewtError("Must specify image to re-sign."))
	}

	imgName := args[0]
	img, err := image.OldImage(imgName)
	if err != nil {
		NewtUsage(nil, err)
		return
	}

	if len(args) > 1 {
		if len(args) > 2 {
			keyId64, err := strconv.ParseUint(args[2], 10, 8)
			if err != nil {
				NewtUsage(cmd,
					util.NewNewtError("Key ID must be between 0-255"))
			}
			keyId = uint8(keyId64)
		}
		keystr = args[1]
		err = img.SetSigningKey(keystr, keyId)
		if err != nil {
			NewtUsage(nil, err)
			return
		}
	}

	err = img.ReSign()
	if err != nil {
		NewtUsage(nil, err)
	}
}

func AddImageCommands(cmd *cobra.Command) {
	createImageHelpText := "Create an image by adding an image header to the " +
		"binary file created for <target-name>. Version number in the header is set " +
		"to be <version>.\n\nTo sign the image give private key as <signing-key> and an optional key-id."
	createImageHelpEx := "  newt create-image my_target1 1.2.0\n"
	createImageHelpEx += "  newt create-image my_target1 1.2.0.3\n"
	createImageHelpEx += "  newt create-image my_target1 1.2.0.3 private.pem\n"
	createImageHelpEx += "  newt create-image my_target1 1.2.0.3 private.pem 5\n"

	createImageCmd := &cobra.Command{
		Use:     "create-image <target-name> <version> [signing-key [key-id]]",
		Short:   "Add image header to target binary",
		Long:    createImageHelpText,
		Example: createImageHelpEx,
		Run:     createImageRunCmd,
	}

	createImageCmd.PersistentFlags().BoolVarP(&newtutil.NewtForce,
		"force", "f", false,
		"Ignore flash overflow errors during image creation")
	createImageCmd.PersistentFlags().BoolVar(&image.UseRsaPss,
		"rsa-pss", false,
		"Use RSA-PSS instead of PKCS#1 v1.5 for RSA sigs")

	cmd.AddCommand(createImageCmd)
	AddTabCompleteFn(createImageCmd, targetList)

	resignImageHelpText := "Sign/Re-sign existing image file. "
	resignImageHelpText += "Note image header will be recreated! "
	resignImageHelpText += "Warning: image has will change if you change key-id "
	resignImageHelpText += "or the type of key used for signing."

	resignImageHelpEx := "  newt resign-image my_target1.img private.pem\n"
	resignImageHelpEx += "  newt resign-image my_target1.img private.pem 5\n"

	resignImageCmd := &cobra.Command{
		Use:     "resign-image <target-name | image-file> [signing-key [key-id]]",
		Short:   "Re-sign image/target using given key.",
		Long:    resignImageHelpText,
		Example: resignImageHelpEx,
		Run:     resignImageRunCmd,
	}

	resignImageCmd.PersistentFlags().BoolVarP(&newtutil.NewtForce,
		"force", "f", false,
		"Ignore flash overflow errors during image creation")

	cmd.AddCommand(resignImageCmd)
}
