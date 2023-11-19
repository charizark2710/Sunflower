import { ErrorMessage } from 'formik';
import './ErrorMessageAtom.atom.scss';

export default function ErrorMessageAtom({name}: {name: string}) {
  return (
    <ErrorMessage className='text-error' name={name} />
  );
}
