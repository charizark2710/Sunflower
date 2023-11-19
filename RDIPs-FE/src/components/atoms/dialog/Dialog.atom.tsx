import { Dialog, DialogContent, DialogTitle } from '@mui/material';

export interface SimpleDialogProps {
  open: boolean;
  title?: string;
  children?: React.ReactNode;
  onClose: (value?: string) => void;
}

export function SimpleDialog(props: SimpleDialogProps) {
  const { onClose, open } = props;

  const handleClose = () => {
    onClose();
  };

  return (
    <Dialog
     onClose={handleClose}
     fullWidth
      open={open}>
      <DialogTitle textAlign={'center'}>{props.title}</DialogTitle>
      <DialogContent>
        {
          props.children
        }
      </DialogContent>
    </Dialog>
  );
}