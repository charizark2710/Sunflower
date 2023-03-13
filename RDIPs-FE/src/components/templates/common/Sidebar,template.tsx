import { Grid } from "@mui/material";

interface SidebarTemplateProps {
    sidebar: React.ReactNode;
}
const SidebarTemplate: React.FC<SidebarTemplateProps> = (props) => {
    return (
        <Grid container spacing={3}>
            <Grid item xs={14}>
                {props.sidebar}
            </Grid>
        </Grid>
    )
};
export default SidebarTemplate;