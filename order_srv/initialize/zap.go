package initialize

import "go.uber.org/zap"

func InitLogger() error {
	logger, err := NewLogger()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}
func NewLogger() (*zap.Logger, error) {
	//cfg := zap.NewProductionConfig()
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"./project.log",
		"stderr",
		"stdout"}
	return cfg.Build()
}
