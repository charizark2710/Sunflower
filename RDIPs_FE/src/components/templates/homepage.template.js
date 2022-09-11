import { Grid } from "@mui/material";

const HomepageTemplate = (props) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                {props.navbar}
            </Grid>
            <div>
                {props.content}
            </div>
            <div>
                {props.footer}
            </div>
        </Grid>
    )
};
export default HomepageTemplate;