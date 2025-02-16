export default () => {
    return require(`${process.cwd()}/config.${process.env.NODE_ENV || 'development'}.json`)
}
