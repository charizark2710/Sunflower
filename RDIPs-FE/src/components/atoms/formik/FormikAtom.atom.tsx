import { Form, Formik } from 'formik';

export interface FormikAtomProps {
  initialValues: object;
  validationSchema?: any;
  children: React.ReactNode;
  onSubmit: (value?: any) => void;
}

export function FormikAtom(props: FormikAtomProps) {
  const { initialValues, validationSchema, children, onSubmit } = props;
  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      <Form>{children}</Form>
    </Formik>
  );
}
