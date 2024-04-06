const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const webpack = require('webpack');

module.exports = {
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
  entry: path.resolve(__dirname, './src/index.tsx'),
  devtool: 'inline-source-map',
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
    }),
    new webpack.DefinePlugin({
      "process.env.REACT_APP_API_URL": JSON.stringify(process.env.REACT_APP_API_URL)
    })
  ],
  stats: 'errors-only',
  devServer: {
    static: [{ directory: path.join(__dirname, "build") }, { directory: path.join(__dirname, "public") }],
    proxy: [{
      context: ['/api'],
      target: {
        host: "localhost",
        protocol: 'http:',
        port: 8080
      },
      pathRewrite: {
        '^/api': ''
      }
    }],
    compress: true,
    historyApiFallback: true,
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