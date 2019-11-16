const path = require("path");
const MiniCSSExtractPlugin = require("mini-css-extract-plugin");
const devMode = process.env.NODE_ENV !== "production";

module.exports = {
  entry: "./src/scripts/index.js",
  mode: process.argv.includes("--release") ? "production" : "development",
  output: {
    filename: "main.js",
    path: path.resolve(__dirname, "dist")
  },
  module: {
    rules: [
      {
        test: /[.](sa|sc|c)ss$/i,
        use: [
          {
            loader: MiniCSSExtractPlugin.loader,
            options: {
              hmr: devMode
            }
          },
          {
            loader: "css-loader",
            options: {
              sourceMap: true
            }
          },
          {
            loader: "sass-loader",
            options: {
              sassOptions: {
                indentWidth: 4
              },
              sourceMap: true
            }
          }
        ]
      }
    ]
  },
  plugins: [
    new MiniCSSExtractPlugin({
      filename: devMode ? "[name].css" : "[name].[hash].css",
      chunkFilename: devMode ? "[id].css" : "[id].[hash].css",
      ignoreOrder: false
    })
  ]
};
