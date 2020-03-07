const path = require('path');
const CssPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const StylelintPlugin = require('stylelint-webpack-plugin');
const devMode = process.env.NODE_ENV !== "production";

/* -----------------------------------
 *
 * CLIENT
 *
 * -------------------------------- */

const client = {
  entry: [
    path.join(__dirname, 'src/styles'),
    path.join(__dirname, 'src/scripts'),
  ],
  mode: process.argv.includes("--release") ? "production" : "development",
  target: 'web',
  context: path.join(__dirname, 'src'),
  output: {
    path: path.resolve(__dirname, "public"),
    filename: "main.js",
  },
  resolve: {
    extensions: [
      '.js',
      '.json',
      '.scss',
      '.ttf',
      '.woff',
      '.svg',
    ],
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  module: {
    rules: [
      {
        enforce: 'pre',
        test: /\.js$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'eslint-loader',
            options: {
              cache: true,
              emitError: true,
              emitWarning: true,
              fix: true,
            }
          }
        ]
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader"
      },
      {
        test: /\.s[ac]ss$/i,
        use: [
          {
            loader: CssPlugin.loader,
            options: {
              hmr: devMode
            }
          },
          {
            loader: 'css-loader',
            options :{
              sourceMap: true
            }
          },
          'postcss-loader',
          {
            loader: 'sass-loader',
            options: {
              sourceMap: true
            },
          },
        ],
      },
    ],
  },
  plugins: [
    new CssPlugin({
      filename: devMode ? "[name].css" : "[name].[hash].css",
      chunkFilename: devMode ? "[id].css" : "[id].[hash].css",
      ignoreOrder: false
    }),
    new CopyPlugin([
      {
        from: './public',
      },
    ]),
    new StylelintPlugin({
      context: path.resolve(__dirname, 'src', 'styles'),
      files: '**/*.s[ac]ss',
      fix: true,
      lintDirtyModulesOnly: true,
    })
  ],
};

/* -----------------------------------
 *
 * EXPORT
 *
 * -------------------------------- */

module.exports = [client];
