# User Service

This service manages user sign-up/sign-in functionality and returns an access token to the client. It integrates seamlessly with AWS Secrets Manager and Cognito for secure authentication.

The service is packaged as a Docker image and pushed to Amazon ECR. A Lambda Function is deployed using the container image. CI/CD is implemented via GitHub Actions, which automatically runs test jobs upon creating a pull request. To deploy the service, trigger the GitHub Action job titled **'Deploy to AWS'**.

## Installation

This service can be run locally using Docker or deployed to AWS using Terraform. Follow the installation steps below based on your preferred environment.

### Mock

[Document](https://vektra.github.io/mockery/latest/)

```cmd
go install github.com/vektra/mockery/v2@v2.50.4
```

### make

- [Download mingw64](https://github.com/niXman/mingw-builds-binaries/releases/download/13.2.0-rt_v11-rev1/x86_64-13.2.0-release-win32-seh-msvcrt-rt_v11-rev1.7z)
- Extract files under `C:\Program Files`
- Add `C:\Program Files\mingw64\bin` to the system's environment variable `Path`
- Rename the executable under the `bin` folder from `mingw32-make` to `make`
- Verify the installation by running the following command in the Command Prompt:

  ```cmd
  make -v
  ```

### Docker

- [WSL](https://learn.microsoft.com/zh-tw/windows/wsl/install)
- [Docker Desktop on Windows](https://docs.docker.com/desktop/install/windows-install/)

### Terraform

- Download [Terraform](https://developer.hashicorp.com/terraform/install#windows)
- Extract the downloaded file to `C:\Program Files\HashiCorp\Terraform`
- Add the Terraform directory to your system's `Path` in the Environment Variables.
- Verify the installation by running the following command in the Command Prompt:

  ```cmd
  terraform version
  ```

### (Optional) AWS

If you want to locally test the Docker image that is built for the Lambda Function.

- [aws-lambda-runtime-interface-emulator](https://github.com/aws/aws-lambda-runtime-interface-emulator/)

## Build Docker Image

```cmd
make build
```

## Run the Docker Image Locally

```cmd
make local_run
```

## Test the API Locally

- Sign Up

```cmd
curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"version":"2.0","path":"/v1/users","httpMethod":"POST","body":"{\"username\":\"<username>\",\"email\":\"<email>\",\"password\":\"<password>\"}","isBase64Encoded":false}'
```

## Deploy the Lambda Function

Use the `make` command to deploy the service to your desired environment:

- Replace `ENV=dev` with your target environment, e.g., dev, stag, prod.
- Ensure that your AWS credentials are properly configured before deployment.

```cmd
make tf_init ENV=dev
```

```cmd
make tf_plan
```

```cmd
make tf_deploy
```
