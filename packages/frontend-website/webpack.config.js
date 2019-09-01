const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const EnvironmentPlugin = require('webpack').EnvironmentPlugin;
const dotenv = require('dotenv');

dotenv.config();

module.exports = {
    mode: 'development',
    entry: './src/bootstrap.tsx',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js'
    },
    devServer: {
        host: "0.0.0.0",
        contentBase: path.join(__dirname, 'dist'),
        compress: true,
        port: 9000,
        historyApiFallback: true
    },
    watchOptions: {

        // watch frontend-shared for changes and recompile
        ignored: [
            /node_modules([\\]+|\/)+(?!@honerlawd\/mentordoc-frontend-shared)/,
            /\@honerlawd\/mentordoc-frontend-shared([\\]+|\/)node_modules/
        ]
    },
    resolve: {
        // Add `.ts` and `.tsx` as a resolvable extension.
        extensions: [".ts", ".tsx", ".js"]
    },
    plugins: [
        new EnvironmentPlugin(Object.keys(process.env)),
        new HtmlWebpackPlugin({
            title: "mentordoc",
            meta: {
                "viewport": "width=device-width, initial-scale=1"
            }
        })
    ],
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                loader: "ts-loader"
            },
            {
                test: /\.scss$/i,
                use: [
                    'style-loader',
                    'css-loader',
                    'sass-loader',
                ],
            },
            {
                test: /\.(svg)$/i,
                use: [
                    {
                        loader: 'file-loader',
                    }
                ],
            },
        ]
    }
};