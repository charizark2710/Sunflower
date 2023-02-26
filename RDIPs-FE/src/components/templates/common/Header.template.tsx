import { Grid } from "@mui/material";

interface HeaderTemplateProps {
    navbar: React.ReactNode;
}
const HeaderTemplate: React.FC<HeaderTemplateProps> = (props) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                {props.navbar}
            </Grid>
        </Grid>
    )
};
export default HeaderTemplate;