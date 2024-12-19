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
