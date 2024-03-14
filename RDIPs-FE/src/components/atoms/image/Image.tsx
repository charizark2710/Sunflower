export interface ImageAtomProps {
  url: string,
  w: string,
  sx?: any,
  alt?: string
}

const Image: React.FC<ImageAtomProps> = ({ url, w }) => {
  return (<img src={url} className="rounded-circle" style={{ width: w }}
    alt="Avatar" />)
}

export default Image;