module.exports = ctx => ({
  map: ctx.env === "development" ? ctx.map : false,
  plugins: {
    stylelint: {
      configFile: ".stylelintrc.json",
      configBaseDir: process.cwd()
    },
    autoprefixer: {},
    cssnano: {}
  }
});
