const path = require('path');
const MiniCss = require('mini-css-extract-plugin');

module.exports = {
    entry: './src/app.js',
    module: {
        rules: [
            {
                test: /\.(js)$/, use: 'babel-loader',
            },
            {
                test: /\.(s*)css$/,
                use: [
                    MiniCss.loader,
                    'css-loader',
                    'postcss-loader',
                    'sass-loader',
                ],
            },
        ],
    },
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'index_bundle.js',
        publicPath: '/',
    },
    target: 'node',
    mode: 'development',
    devServer: {
        publicPath: '/',
        contentBase: './src',
        hot: true,
        port: 4000,
        historyApiFallback: true,
    },
    plugins: [
        new MiniCss({
            filename: 'style.css',
        }),
    ],
    resolve: {
        alias: {
            utils: path.resolve(__dirname, 'src/public/js/utils'),
            components: path.resolve(__dirname, 'src/public/js/components'),
            partials: path.resolve(__dirname, 'src/public/js/partials'),
        },
    },
};
