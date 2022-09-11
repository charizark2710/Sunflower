import { Grid } from "@mui/material";

const HeaderTemplate = (props) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                {props.navbar}
            </Grid>
        </Grid>
    )
};
export default HeaderTemplate;