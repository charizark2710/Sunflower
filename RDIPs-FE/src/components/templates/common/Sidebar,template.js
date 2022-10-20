import { Grid } from "@mui/material";

const SidebarTemplate = (props) => {
    return (
        <Grid container spacing={3}>
            <Grid item xs={14}>
                {props.sidebar}
            </Grid>
        </Grid>
    )
};
export default SidebarTemplate;