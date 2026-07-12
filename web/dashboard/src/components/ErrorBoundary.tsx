import { Component, type ComponentChild } from 'preact';

interface Props {
  children: ComponentChild;
  fallback?: ComponentChild;
}

interface State {
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        this.props.fallback ?? (
          <div class="error-banner" role="alert">
            <strong>Something went wrong.</strong>{' '}
            {this.state.error.message}
          </div>
        )
      );
    }
    return this.props.children;
  }
}
