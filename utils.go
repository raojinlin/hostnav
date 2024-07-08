package jmfzf

import "encoding/json"

func MapToStruct[T CloudProviderConfig | JumpServerCofnig | SshConfig | DockerConfig | KubernetesConfig](m interface{}, v *T) error {
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
