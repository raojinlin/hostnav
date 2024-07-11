package hostnav

import "encoding/json"

func MapToStruct[T CloudProviderOption | JumpServerOption | SshConfig | DockerOption | KubernetesOption | FileOption](m interface{}, v *T) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
