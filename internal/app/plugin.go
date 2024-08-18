package app

import (
	"github.com/Yuelioi/vidor/internal/shared"
)

func loadLocalPlugin(pluginPath string) (shared.Downloader, error) {
	// client := plugin.NewClient(&plugin.ClientConfig{
	// 	HandshakeConfig: shared.HandshakeConfig,
	// 	Plugins: map[string]plugin.Plugin{
	// 		"downloader": &shared.DownloaderRPCPlugin{},
	// 	},
	// 	Cmd: exec.Command(pluginPath),
	// })

	// rpcClient, err := client.Client()
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating client for plugin %s: %v", pluginPath, err)
	// }

	// raw, err := rpcClient.Dispense("downloader")
	// if err != nil {
	// 	return nil, fmt.Errorf("error dispensing plugin %s: %v", pluginPath, err)
	// }

	// downloader, ok := raw.(shared.Downloader)
	// if !ok {
	// 	return nil, fmt.Errorf("plugin %s does not implement the expected interface", pluginPath)
	// }

	// return downloader, nil
	return nil, nil
}
