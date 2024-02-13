# Environment Variables

## Version

Version info

 - `SOURCE_BRANCH` - 
 - `SOURCE_COMMIT` - 
 - `IMAGE_NAME` - 

## Config

Application configuration

 - `ENV` (default: `production`) - Env must be local, development, test or production
 - `HOST` - 
 - `PORT` (default: `3000`) - 
 - `API_KEY` (**required**, non-empty) - 
 - Repository configuration
   - `REPO_CONN` (**required**, non-empty) - Database connection string
   - `REPO_DEACTIVATION_PERIOD` (default: `8h`) - 
   - `REPO_CONN` (**required**, non-empty) - Connection string
   - `REPO_NAME` (**required**, non-empty) - Index Name
   - `REPO_RETENTION` (default: `5`) - Index retention
 - Repository configuration
   - `INDEX_CONN` (**required**, non-empty) - Database connection string
   - `INDEX_DEACTIVATION_PERIOD` (default: `8h`) - 
   - `INDEX_CONN` (**required**, non-empty) - Connection string
   - `INDEX_NAME` (**required**, non-empty) - Index Name
   - `INDEX_RETENTION` (default: `5`) - Index retention
