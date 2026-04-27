# Sunflower Deployment

All deployment scripts and compose files are now located in this `deployment` folder.

## Usage

- Use the scripts and compose files here for all deployment operations.
- Make sure to update environment files (`local.deploy.env`, `local.env`, `local_ssl.env`) in this folder as needed.
- All relative paths in compose files are updated to work from this folder.

## Example

```bash
cd deployment
bash deploy-swarm.sh
```
