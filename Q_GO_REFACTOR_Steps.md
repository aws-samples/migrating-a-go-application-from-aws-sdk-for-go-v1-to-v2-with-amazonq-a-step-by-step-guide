# Migrating a Go application from AWS SDK Go V1 to Go V2 with Amazon Q Developer

## Prerequisites 

For this part of the blog you should have the following prerequisites: 

* An AWS [account](https://signin.aws.amazon.com/signin?redirect_uri=https%3A%2F%2Fportal.aws.amazon.com%2Fbilling%2Fsignup%2Fresume&client_id=signup)

* [Visual Studio](https://aws.amazon.com/visualstudiocode/) Code or [Jet Brains IDE](https://docs.aws.amazon.com/toolkit-for-jetbrains/latest/userguide/setup-toolkit.html) with AWS Toolkit Installed

    * Installation instructions for AWS Toolkit can be found in the Getting Started for [VS Code](https://docs.aws.amazon.com/toolkit-for-vscode/latest/userguide/setup-toolkit.html) and  [JetBrains](https://docs.aws.amazon.com/toolkit-for-jetbrains/latest/userguide/setup-toolkit.html)

* For Visual Studio IDE install the [Amazon Q from the Visual Studio Marketplace](https://marketplace.visualstudio.com/items?itemName=AmazonWebServices.amazon-q-vscode) 

* For JetBrains IDE install the [Anazon Q from the JetBrains Markerplace](https://plugins.jetbrains.com/plugin/24267-amazon-q/)

    * For Both [Visual Studio](https://marketplace.visualstudio.com/items?itemName=AmazonWebServices.amazon-q-vscode) and [JetBrains](https://plugins.jetbrains.com/plugin/24267-amazon-q/) IDEs you will need to Authenticate using free builder ID or IAM Identity Center with a Amazon Q Developer Pro subscription as documented in [Authenticating in Visual Studio Code](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/q-in-IDE-setup.html#setup-vscode) or [Authenticating in JetBrains IDEs](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/q-in-IDE-setup.html#setup-jetbrains)

## Analyze Current Code with Amazon Q Developer
The Amazon Q Developer Agent for software development can help you develop code features or make code changes to projects in your integrated development environment (IDE). You explain the task you want to accomplish, and Amazon Q uses the context of your current project or workspace to generate code to implement the changes. Amazon Q can help you build AWS projects or your own applications.

The following steps will guide you through using Amazon Q developer in our IDE to help us get familiar our Application code using AWS SDK for GO v1 and to create a "CodeExplaination.md" file using [/dev](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/software-dev.html#develop) 

1. Open the git clone project folder ***[Go SDK Amazon Sample Project](https://gitlab.aws.dev/go_sdk_blog/go_sdk_amazon_q)*** < **Replace with Public Link**> from the previous steps in your IDE.
2. Click the Amazon Q Developer icon in your IDE to show the Amazon Q chat panel "<img align="center" width="25" height="25" src="./images/amazon_q.png">"

2. Enter /dev ([this alert Q that we want to do development code action](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/software-dev.html)) in the Amazon Q chat panel followed by a description of the task you want to accomplish or the issue you want to resolve. Use the prompt below:

```text
/dev
"Generate a CodeExplaination.md for the Go project. Include sections for:

1. Project title and brief description
2. Installation instructions
3. Basic usage example
4. Key features
5. Dependencies
6. How to contribute
7. License information

Keep the content clear, concise, and informative. 
Use proper Markdown formatting. 
Assume the project is open-source and hosted on GitHub. 
Tailor the content to be relevant for Go developers. 
Limit the total length to approximately 500 words."
```




