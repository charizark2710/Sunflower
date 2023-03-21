import { Grid } from "@mui/material";
interface HeaderTemplateProps {
    header: React.ReactNode;
}
const HeaderTemplate: React.FC<HeaderTemplateProps> = (props) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                {props.header}
            </Grid>
        </Grid>
    )
};
export default HeaderTemplate;