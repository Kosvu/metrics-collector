package agent

func isRetriableNet(err error) bool {
	return err != nil
}
