import * as React from 'react'
import { Button, Alert, Intent, IButtonProps } from '@blueprintjs/core'

type OwnProps = {
    confirmationText: string,
    onDestroy: () => void
}

type State = {
    isConfirming: boolean
}

type Props = OwnProps & IButtonProps

const { Fragment } = React

class DestroyButton extends React.Component<Props, State> {
    constructor(props: OwnProps) {
        super(props)
        this.state = {
            isConfirming: false
        }
    }

    render() {
        const { text, confirmationText, minimal, onDestroy } = this.props

        return (
            <Fragment>
                <Button
                    text={ text }
                    minimal={ minimal }
                    icon="trash"
                    intent={ Intent.DANGER }
                    onClick={ (e: React.SyntheticEvent<HTMLElement>) => {
                        e.stopPropagation()
                        this.setState({ isConfirming: true })
                    } }
                />
                { this.state.isConfirming && (
                    <Alert
                        icon="trash"
                        intent={ Intent.DANGER }
                        isOpen={ this.state.isConfirming }
                        onConfirm={ (e: React.SyntheticEvent<HTMLElement>) => {
                            e.stopPropagation()
                            this.setState({ isConfirming: false })
                            onDestroy()
                        } }
                        onCancel={ (e: React.SyntheticEvent<HTMLElement>) => {
                            e.stopPropagation()
                            this.setState({ isConfirming: false })
                        } }
                        cancelButtonText="Cancel"
                        canEscapeKeyCancel
                        canOutsideClickCancel
                    >
                        { confirmationText }
                    </Alert>
                ) }
            </Fragment>
        )
    }
}

export default DestroyButton
