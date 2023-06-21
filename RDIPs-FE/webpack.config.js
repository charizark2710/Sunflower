const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  entry: path.resolve(__dirname, './src/index.tsx'),
  output: {
    path: path.join(__dirname, "/build"),
    filename: "bundle.js",
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.scss', '.css', ".json"],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: "public/index.html",
      filename: 'index.html',
      favicon: 'public/favicon.ico',
      manifest: "public/manifest.json"
    }),
  ],
  stats: 'errors-only',
  devServer: {
    static: path.join(__dirname, "build"),
    compress: true,
    port: 3030, // you can change the port
  },
  module: {
    rules: [
      {
        test: /\.(ts|js)x?$/, // .js and .jsx files
        exclude: /node_modules/, // excluding the node_modules folder
        use: [
          {
            loader: "babel-loader",
          },
          {
            loader: "ts-loader",
          }
        ],
      },
      {
        test: /\.(sa|sc|c)ss$/, // styles files
        use: ["style-loader", "css-loader", "sass-loader"],
      },
      {
        test: /\.json$/,
        use: ['json-loader'],
        type: 'javascript/auto'
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg|jpg|jpeg|gif|ico)$/, // to import images and fonts
        exclude: /node_modules/,
        use: ['file-loader?name=[name].[ext]']
      },
    ],
  },
};