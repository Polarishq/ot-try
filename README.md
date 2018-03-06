# ot-try

## Setup ECS Registry
* [Create a new repository](https://us-west-1.console.aws.amazon.com/ecs/home?region=us-west-1#/repositories/create/new) with the name ot-try

## Setup Codeship SSH Key in GitHub

### Add ot-try to the Machine Users - READONLY group in GitHub
* Navigate to the ["Collaborators & teams" section](https://github.com/Polarishq/ot-try/settings/collaboration) of your repository's settings page
* Under the "Teams" section, click the dropdown in the bottom left to "Add a team"
* Select "Machine Users - READONLY"
* Your repository should show up in the list of repositories for the ["Machine Users - READONLY" group in GitHub](https://github.com/orgs/Polarishq/teams/machine-users-readonly/repositories)

### Create Codeship project
* Open Codeship Web UI in the [Nova organization](https://app.codeship.com/nova)
* Go to projects and click on "New project"
* Connect to Github repository
* Paste the git clone URL for your repo `git@github.com:Polarishq/ot-try.git`

### Remove the CodeShip public SSH key from your repo's deploy keys
* Navigate to your repo's [GitHub Project Settings](https://github.com/Polarishq/ot-try/settings/keys) 
* Click the "Delete" button for the "codeship" key
* Confirm the delete action in the popup dialog box

### Add the SSH Public key to the polaris-codeship GitHub Machine User Account
* Open Codeship Web UI in the [Nova organization](https://app.codeship.com/nova)
* Click on "Select Project"
* Click on the settings icon for your project
* Copy the "SSH public key" to the clipboard
* Use 1Password to log into GitHub using the "polaris-codeship" machine user credentials (The 2 factor auth code is also setup in 1Password)
* Go to the [polaris-codeship GitHub user settings page](https://github.com/settings/keys)
* Click on "SSH and GPG Keys" in the left navbar
* Click "New SSH key"
* Enter the following
	* Title: codeship-ot-try
	* Key: The public SSH key you copied from CodeShip
* Click "Add SSH Key"

## Setup Codeship badge
* Open Codeship Web UI in the [Nova organization](https://app.codeship.com/nova)
* Click on "Select Project"
* Click on the settings icon for your project
* Go to General Settings
* Scroll down to "Status Images"
* Click "Copy Markdown Syntax"
* Paste the Markdown into this README.md in place of these instructions
* Commit your changes
    * git commit -am "Setup codeship badge"
* Push your branch up to github
    * git push origin master

## Setup Coveralls project
* Join the coveralls PolarisHQ organization on https://coveralls.io/
* Navigate to https://coveralls.io/repos/new and enable this repo
* Go to the settings for the coveralls project and copy the "REPO TOKEN"
* Paste the repo token into credentials.secret.env
* Run "make encrypt"
* Commit the encrypted file credentials.env.encrypted to your github repo
* Create and checkout a new branch (not master)
* Run "make codeship"

### Setup minimum code coverage requirements
* Navigate to https://coveralls.io/github/Polarishq/ot-try/settings
* Setup minimum code coverage requirements
    * COVERAGE THRESHOLD FOR FAILURE                90%
    * COVERAGE DECREASE THRESHOLD FOR FAILURE       5%

### Setup Coveralls badge
* Navigate to https://coveralls.io/github/Polarishq/ot-try/settings
* Add the coveralls badge to README.md in place of these instructions
    * Click "Embed"
    * Copy Markdown to README.md
* Commit your changes
    * git commit -am "Setup coveralls"
* Push your branch up to github
    * git push origin master

## Setup Github branch protection
* Navigate to the branch settings of your repo in github [your repo on gihub](https://github.com/Polarishq/ot-try/settings/branches)
* Add "master" to the "Protected branches"
    * Enable "Protect this branch"
    * Enable "Require status checks to pass before merging"
        * Enable "Include administrators"
        * Enable "Require branches to be up to date before merging"
    * Enable status checks
        * continuous-integration/codeship
        * continuous-integration/coveralls
        
