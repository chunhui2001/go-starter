## 项目主页
https://console.cloud.google.com/

### OAuth 2.0 用戶端 ID
## Creating client IDs
1. In the Google Cloud console, go to the Credentials page.
> https://console.cloud.google.com/apis/credentials

2. From the projects list, select the project containing your API.
3. If this is your first time creating a client ID in this project, use the sub-steps to go to the OAuth consent page; otherwise, skip to the next step.
	-- Click OAuth consent screen.
	-- Enter a name in the Application name field.
	-- Fill out the rest of the fields as needed.
	-- Click Save.
4. In the Create credentials drop-down list, select OAuth client ID.
5. Select Web application as the application type.
6. In Name, enter a name for your client ID.
7. In Authorized JavaScript origins, enter one of the following:
	-- http://localhost:8080 if you are testing the backend locally.
	-- https://YOUR_PROJECT_ID.appspot.com, replacing YOUR_PROJECT_ID with your App Engine project ID if you are deploying your backend API to your production App Engine.

	-- Important: You must specify the site or host name under Authorized JavaScript origins using https for this to work with Cloud Endpoints, unless you are testing with localhost, in which case you must use http.

8. 	Click Create.
	-- You use the generated client ID in your API backend and in your client application.


