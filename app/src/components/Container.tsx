import * as React from 'react'

import styles from './Container.module.scss'

const Container: React.SFC<{}> = ({ children }) => {
    return (
        <div className={ styles.container }>
            { children }
        </div>
    )
}

export default Container
