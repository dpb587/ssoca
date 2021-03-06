## Usage Details

To create and automatically install a Tunnelblick profile you may use the `--install` option...

    $ ssoca openvpn create-tunnelblick-profile --install
    The profile 'acme-prod-aws-use1' has successfully been installed.

To create only the profile files, leave off `--install` (it will create a `{service-name}.tblk` in your current directory)...

    $ ssoca openvpn create-tunnelblick-profile

Once you are ready to install the profile, open the file to register it with Tunnelblick...

    $ open openvpn.tblk

Tunnelblick will give two prompts regarding the security concern of profiles with custom scripts, followed by a prompt about which context to install the profile, and finally an authentication prompt to install the profile. Once installed, you can connect/disconnect like other OpenVPN profiles.


## Notes

The generated profile includes scripts to automatically regenerate connection certificates upon establishing a VPN connection. There are a few caveats...

 * The scripts execute as `root`, therefore Tunnelblick gives very clear warnings and confirmations when first installing the profile. It is referring to the `pre-connect.sh` file inside the generated `*.tblk` directory.
 * For security, Tunnelblick may keep a shadow copy of profile configuration which requires a user authentication prompt whenever the configuration is changed. The ssoca profiles have a short lifetime and are regenerated every connection. To avoid this user intervention, the background scripts automatically updates the shadow copy.
 * If your authentication session expires, Tunnelblick may error during the next connection attempt with an [ambiguous] error dialog.
 * When the connection is lost, the `openvpn` process will exit and, by default, not attempt to reconnect. To change this behavior, navigate to the Advanced settings of your profile and enable the "Keep connected" checkbox of the "Connecting & Disconnecting" tab.
