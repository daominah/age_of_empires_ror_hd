# GitHub Actions Workflow Setup

> **Note:** You are working with two GitHub accounts:
> - **daominah/age_of_empires_ror_hd** (this repository, where the workflow runs)
> - **aoe2rating/aoe2rating.github.io** (the deployment target, where the generated site is published)
>
> Make sure you are logged into the correct account for each step. If you need to switch accounts, use an incognito/private browser window or log out and log in as needed.

- [daominah/age_of_empires_ror_hd](https://github.com/daominah/age_of_empires_ror_hd)
- [aoe2rating/aoe2rating.github.io](https://github.com/aoe2rating/aoe2rating.github.io)

## Daily AoE2 Rating Update Workflow

This workflow automatically:

1. Runs the AoE2 rating percentile generator daily at 12:00 ICT
2. Generates `index.html` and `output_charts/` directory
3. Deploys these files to `aoe2rating/aoe2rating.github.io` repository

## Setup Instructions

### 1. Create a Fine-Grained Personal Access Token (PAT) in the `aoe2rating` Account

You need to create a GitHub Fine-Grained Personal Access Token with **write access specifically to the `aoe2rating/aoe2rating.github.io` repository** (log in as `aoe2rating`):

1. **Log in** to your **`aoe2rating` account**.
2. [Create a new fine-grained token](
   https://github.com/settings/personal-access-tokens/new)
3. Token name: `aoe2rating_github_io_deploy`
4. Description: `Token for running workflow from daominah/age_of_empires_ror_hd to deploy to aoe2rating.github.io`
5. Resource owner: keep default selected `aoe2rating`
  - This means the token can only write/modify resources owned by `aoe2rating`
  - The token can read all public repositories, but can only **write** to `aoe2rating`'s repositories
  - Since `aoe2rating.github.io` is owned by `aoe2rating`, this is correct
6. Expiration: No expiration.
7. Repository access:
   Select "Only select repositories" and
   choose `aoe2rating/aoe2rating.github.io`
8. Permissions:
   `Repositories` > `Add permission` > `Contents` > `Access: Read and write`
9. Generate token
10. **Copy the token immediately** (you won't be able to see it again)

### 2. Add the Token as a Repository Secret in `daominah/age_of_empires_ror_hd`

1. Log in to your `daominah` account (if needed).
2. Go to your repository's [Actions secrets page](
   https://github.com/daominah/age_of_empires_ror_hd/settings/secrets/actions)
3. Click "New repository secret"
4. Name: `GH_PAGES_AOE2RATING_TOKEN`
5. Value: Paste the token you copied from the `aoe2rating` account
6. Check a new row "GH_PAGES_AOE2RATING_TOKEN" is added to "Repository secrets".

### 3. Verify the Workflow in `daominah/age_of_empires_ror_hd`

1. Go to the [Actions tab](
   https://github.com/daominah/age_of_empires_ror_hd/actions) in your repository
2. You should see "Daily AoE2 Rating Update" workflow
3. You can manually trigger it using "Run workflow" button to test
4. The workflow will run automatically every day at 12:00 ICT

## Manual Trigger

You can manually trigger the workflow at any time:

- Go to [Actions tab](https://github.com/daominah/age_of_empires_ror_hd/actions) → "Daily AoE2 Rating Update" → "Run workflow"

## Troubleshooting

- If the workflow fails with authentication errors, verify that `GH_PAGES_AOE2RATING_TOKEN` secret is set correctly in `daominah/age_of_empires_ror_hd` and that the token has **push (write) access** to `aoe2rating/aoe2rating.github.io`.
- If the workflow fails with "repository not found", ensure the token has access to `aoe2rating/aoe2rating.github.io`.
- Check the Actions logs for detailed error messages
