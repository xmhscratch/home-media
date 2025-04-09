export default () => {
    const envName: string = process.env['NODE_ENV'] || 'development';

    return {
        env: envName,
        ...require(`${process.cwd()}/config.${envName}.json`),
    }
}
