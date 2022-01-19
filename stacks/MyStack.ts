import * as sst from "@serverless-stack/resources";

export default class MyStack extends sst.Stack {
  constructor(scope: sst.App, id: string, props?: sst.StackProps) {
    super(scope, id, props);

    // Create a HTTP API
    const api = new sst.Api(this, "Api", {
      routes: {
        "POST /": "src/main.go",
      },
    });

    api.attachPermissions(["ses"]);

    this.addOutputs({
      ApiEndpoint: api.url,
    });
  }
}
