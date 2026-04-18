@echo off
rem Convenience wrapper: run the fabsim CLI from the repo root via `go run`.
rem %~dp0 expands to the directory this .bat lives in, so invocation from a
rem subdirectory still resolves the cmd/fabsim package correctly.
go run "%~dp0cmd\fabsim" %*
